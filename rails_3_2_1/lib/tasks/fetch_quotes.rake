namespace :fetch_quotes do
  desc "fetch stock quotes from yahoo by days"
  task :fetch_days, [:stock_code, :days] => :environment do |t, args|
    s = StockMap.find_by_stock_id(args[:stock_code])
    s.fetch_history(args[:days].to_i)
  end
end
