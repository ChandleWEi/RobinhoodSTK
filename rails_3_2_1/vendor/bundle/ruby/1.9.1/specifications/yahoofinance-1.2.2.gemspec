# -*- encoding: utf-8 -*-

Gem::Specification.new do |s|
  s.name = "yahoofinance"
  s.version = "1.2.2"

  s.required_rubygems_version = nil if s.respond_to? :required_rubygems_version=
  s.authors = ["Nicholas Rahn"]
  s.autorequire = "yahoofinance"
  s.cert_chain = nil
  s.date = "2007-09-06"
  s.email = "nick @nospam@ transparentech.com"
  s.require_paths = ["lib"]
  s.required_ruby_version = Gem::Requirement.new("> 0.0.0")
  s.rubygems_version = "1.8.23"
  s.summary = "A package for retrieving stock quote information from Yahoo! Finance."

  if s.respond_to? :specification_version then
    s.specification_version = 1

    if Gem::Version.new(Gem::VERSION) >= Gem::Version.new('1.2.0') then
    else
    end
  else
  end
end
