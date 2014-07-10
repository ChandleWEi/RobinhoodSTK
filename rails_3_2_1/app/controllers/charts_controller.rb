class ChartsController < ApplicationController
  def index
    days = 35
    id = 1
    #TODO bugs here. if days is longer than smas26.length
    #smas12 = StockMap.find(id).smas.where("duration = ?", 12).limit(days).reverse
    #smas26 = StockMap.find(id).smas.where("duration = ?", 26).limit(days).reverse
    emas12 = StockMap.find(id).emas.where("duration = ?", 12).limit(days).reverse
    emas26 = StockMap.find(id).emas.where("duration = ?", 26).limit(days).reverse
    close_prices = StockMap.find(id).rawdata.limit(days).reverse
    @emas = LazyHighCharts::HighChart.new('graph') do |f|
      f.title({ :text=>"chart"})
      f.options[:xAxis][:categories] = close_prices.map{|x|x.record_date}
      #f.labels(:items=>[:html=>"moving averages", :style=>{:left=>"40px", :top=>"8px", :color=>"black"} ])      
      #f.series(:type=> 'column',:name=> 'Jane',:data=> [3, 2, 1, 3, 4])
      #f.series(:type=> 'spline',:name=> 'Average', :data=> [3, 2.67, 3, 6.33, 3.33])
      f.series(:type=> 'spline',:name=> 'close_price', :data=> close_prices.map{|x|x.adjusted_close_price.to_f})
      #f.series(:type=> 'spline',:name=> 'sma(12)', :data=> smas12.map{|x|x.sma.to_f})
      #f.series(:type=> 'spline',:name=> 'sma(26)', :data=> smas26.map{|x|x.sma.to_f})
      f.series(:type=> 'spline',:name=> 'ema(12)', :data=> emas12.map{|x|x.ema.to_f})
      f.series(:type=> 'spline',:name=> 'ema(26)', :data=> emas26.map{|x|x.ema.to_f})
      #f.series(:type=> 'pie',:name=> 'Total consumption', 
      #         :data=> [
      #           {:name=> 'Jane', :y=> 13, :color=> 'red'}, 
      #           {:name=> 'John', :y=> 23,:color=> 'green'},
      #           {:name=> 'Joe', :y=> 19,:color=> 'blue'}
      #],
      #  :center=> [100, 80], :size=> 100, :showInLegend=> false)
    end

    difs = StockMap.find(id).difs.limit(days).reverse
    macds = StockMap.find(id).macds.limit(days).reverse
    @macds = LazyHighCharts::HighChart.new('graph') do |f|
      f.title({ :text=>"chart"})
      f.options[:xAxis][:categories] = difs.map{|x|x.record_date}
      f.series(:type=> 'spline',:name=> 'dif', :data=> difs.map{|x|x.dif.to_f})
      f.series(:type=> 'spline',:name=> 'macd', :data=> macds.map{|x|x.macd.to_f})
    end

    respond_to do |f|
      f.html
      f.json { render :json => @emas }
    end
  end
end
