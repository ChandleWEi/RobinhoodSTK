class StockMap < ActiveRecord::Base
  attr_accessible :stock_id, :stock_name
  has_many :rawdata, :conditions => "volume > 0", :order => "record_date desc"
  has_many :smas, :order => "record_date desc" do
    def duration(n)
      self.where(:duration => n)
    end
  end
  has_many :emas, :order => "record_date desc"  do
    def duration(n)
      self.where(:duration => n)
    end
  end
  has_many :difs, :order => "record_date desc"
  has_many :macds, :order => "record_date desc"

  #fetch and store stock historical rawdata
  def fetch_history(days)
    raw = YahooFinance::get_historical_quotes_days(self.stock_id, days)
    raw.each do |r|
      rd = Rawdatum.find_or_create_by_stock_map_id_and_record_date(self.id, r[0])
      rd.open_price = r[1]
      rd.high_price = r[2]
      rd.low_price = r[3]
      rd.close_price = r[4]
      rd.volume = r[5]
      rd.adjusted_close_price = r[6]
      rd.save
    end
  end

  #initialize sma
  def init_sma(duration)
    all_data = self.rawdata.reverse
    (all_data.length-duration+1).times do |d|
      res = (all_data[d..(d+duration-1)].inject(0){|a,b|a+=b.adjusted_close_price}/duration).to_f
      self.smas.create(:duration => duration, :record_date => all_data[d+duration-1].record_date, :sma => res)
    end
  end

  def generate_sma(duration)
    old_sma = self.smas.duration(duration).limit(duration-1).reverse
    if old_sma.empty?
      self.init_sma(duration)
    else
      new_data = self.rawdata.where("record_date > ?", old_sma.last.record_date)
      all_data = (new_data + self.rawdata.where("record_date <= ?", old_sma.last.record_date).limit(duration-1)).reverse
      new_data.reverse.each_with_index do |nd,i|
        new_sma = (old_sma[i].sma + (nd.adjusted_close_price - all_data[i].adjusted_close_price)/duration).to_f
        self.smas.create(:duration => duration, :record_date => nd.record_date, :sma => new_sma)
      end
    end
  end

  #generate ema, using the latest ema or the earliest sma
  def generate_ema(duration)
    start_ema = 0
    sma = 0
    old_emas = self.emas.duration(duration)
    if old_emas.empty?
      #use the earliest sma as original ema
      if self.smas.duration(duration).empty?
        self.init_sma(duration)
      end
      sma = self.smas.duration(duration).reverse.first
      if sma.nil?
        p "not enough rawdata to initialize sma. neither nor ema."
        return false
      else
        #Ema.find_or_create_by_stock_map_id_and_duration_and_record_date(self.id, duration, sma.sma)
        start_ema = self.emas.create(:duration => duration, :ema => sma.sma, :record_date => sma.record_date)
      end
    else
      start_ema = old_emas.first
    end
    records = self.rawdata.where("record_date > ?", start_ema.record_date).reverse
    exp = 2/(duration + 1)
    res = start_ema.ema.to_f
    records.each do |r|
      #ema = Ema.find_or_create_by_stock_map_id_and_duration_and_record_date(self.id, duration, records[n].record_date)
      new_ema = self.emas.create(:duration => duration, :record_date => r.record_date)
      res = (exp * r.adjusted_close_price + (1-exp) * res).to_f
      new_ema.ema = res
      new_ema.save
    end
  end

  #initialize dif
  def generate_dif(d1, d2)
    if d1 >= d2
      p "wrong dates"
      return false
    end
    self.generate_sma(d1)
    self.generate_sma(d2)
    self.generate_ema(d1)
    self.generate_ema(d2)
    ema1 = self.emas.duration(d1)
    ema2 = self.emas.duration(d2)
    if ema2.empty? || ema1.empty?
      p "some emas not exist"
      return false
    else
      start_date = self.difs.empty? ? ema2.last.record_date : self.difs.first.record_date
      ema2.reverse!
      ema1.select{|x|x.record_date>=start_date}.reverse.each_with_index do |e1,i|
        self.difs.create(:record_date => e1.record_date, :dif => (e1.ema-ema2[i].ema).to_f)
      end
    end
  end

  #initialize macd
  def generate_macd
    d = 9
    exp = 2/(d+1)
    old_macd = self.macds.first
    if old_macd.nil?
      return false if self.difs.size < d
      start_date = self.difs[-d].record_date
      start_macd = (self.difs[-d..-1].inject(0){|x,y|x+=y.dif}/d).to_f
      old_macd = self.macds.create(:record_date => start_date, :macd => start_macd)
    end
    start_date = old_macd.record_date
    self.difs.where("record_date > ?", start_date).reverse.each do |dif|
      macd = old_macd.macd
      res = (exp * dif.dif + (1-exp) * macd).to_f
      old_macd = self.macds.create(:record_date => dif.record_date, :macd => res)
    end
  end
  #EMA(today) = exp*close_price(today) + (1-exp)*EMA(yesterday) or exp*(close_price(today) - EMA(yesterday))+EMA(yesterday)
end
