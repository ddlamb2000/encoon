class DefaultDatesInTables < ActiveRecord::Migration
  def self.up
    Workspace.find(:all).each do |item|
      puts "Update workspace " + item.id.to_s
      item.begin = Date::civil(1,1,1)
      item.end = Date::civil(9999,12,31)
      item.save! if item.valid?
    end

    Grid.find(:all).each do |item|
      puts "Update grid " + item.id.to_s
      item.begin = Date::civil(1,1,1)
      item.end = Date::civil(9999,12,31)
      item.save! if item.valid?
    end

    Column.find(:all).each do |item|
      puts "Update column " + item.id.to_s
      item.begin = Date::civil(1,1,1)
      item.end = Date::civil(9999,12,31)
      item.save! if item.valid?
    end

    Row.find(:all).each do |item|
      puts "Update row " + item.id.to_s
      item.begin = Date::civil(1,1,1)
      item.end = Date::civil(9999,12,31)
      item.save! if item.valid?
    end

    User.find(:all).each do |item|
      puts "Update user " + item.id.to_s
      item.begin = Date::civil(1,1,1)
      item.end = Date::civil(9999,12,31)
      item.save! if item.valid?
    end
  end

  def self.down
  end
end
