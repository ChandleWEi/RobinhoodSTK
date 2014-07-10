#!/System/Library/Frameworks/Ruby.framework/Versions/1.8/usr/bin/ruby -W0
# -*- coding: utf-8 -*-
require 'iconv'
require 'open-uri'
require 'timeout'

# 上证指数,2577.429 （今日开盘）,2559.911（昨日收盘）,2555.284当前,2590.776今日最高,2554.787最低,0,0,92035473手成交量,89024533620成交额,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2009-05-05,11:35:58

# http://image.sinajs.cn/newchart/macd/n/sh600705.gif
# stock_str = sotckids.split('').join(',')
class Stock

  HQ_URL = "http://hq.sinajs.cn/list=%s"
  STOCK_IDS = "sh000001,sh600466,sz002109,sh600255,sz002124,sz000883,sh600031,sz002020,sz002444,sz000727,sh600482,sh600283,sz000998,sh600897,sz001896,sh600636,sz002083,sh600466,sh600073,sh600702,sz000710,sz002143,sh600850,sh600757,sz002312,sh600369,sh600036,sh600000,sh600581,sh601099,sh600369,sh600729,sz000863,sh600116,sh600503,sh600757,sh600251,sh000001,sh600705"
  #STOCK_IDS = "sh600705"


  def self.get_stock
    times = 999
    time = 10
    begin
      timeout(time) do 
        url = HQ_URL % STOCK_IDS
        #p url
        data = open(url)
        stock_temp_data = Iconv.iconv('utf-8', 'gbk', data.read)
        stock_data_array = stock_temp_data[0].split(/\n/)
        stock_data_array.each do |stock_data|
          data_entry = stock_data[/="(.*)"/,1]
          data_name =  stock_data[/hq_str_(.*)=/,1]
          data_array = data_entry.split( /,/)
          change_percent = ((data_array[3].to_f/data_array[2].to_f) - 1) * 100
          data_average = (data_array[4].to_f + data_array[5].to_f)/2
          printf("%s(%s)[%0.2f/%+0.2f]%0.2f|%0.2f|%0.2f(%s|%s)<%s>\n", data_name,data_array[0],data_array[3],change_percent,data_array[4],data_array[5],data_average,data_array[8], data_array[9],data_array[31])

          #        num = 10.25
          #        if data_name =~/600036/
          #          if data_array[3].to_f > num
          #            system "/usr/bin/terminal-notifier -message \"#{data_name}|#{data_array[0]}|#{data_array[3]}\""
          #          end
          #        end
        end
      end
    rescue Timeout::Error => e
      times -= 1
      if times > 0
        puts "----->try times left #{times}"
        sleep 0.42 and retry
      end      
    rescue => e
      p "error is #{e}"
    end
  end
end


while(1) do
  sleep 3
  Stock.get_stock()
end

