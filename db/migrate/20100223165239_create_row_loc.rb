class CreateRowLoc < ActiveRecord::Migration
  def self.up
    create_table :row_locs do |t|
      t.string   :uuid, :limit => 36
      t.integer  :version
      t.string   :locale, :limit => 10
      t.string   :base_locale, :limit => 10
      t.string   :name
      t.text     :description
      t.integer  :lock_version, :default => 0
    end

    add_index :row_locs, [:uuid, :version, :locale]

    Row.find(:all).each do |item|
      puts "Update row " + item.id.to_s
      LANGUAGES.each do |lang, locale|
        puts "Update locale " + locale.to_s
        loc = RowLoc.new
        loc.uuid = item.uuid
        loc.version = item.version
        loc.locale = locale
        loc.base_locale = I18n.default_locale.to_s
        loc.name = item.name + (locale != "en" ? " (" + locale + ")" : "")
        loc.description = item.description + (locale != "en" ? " (" + locale + ")" : "") if item.description.present?
        loc.save!
      end
    end
  end

  def self.down
    drop_table :row_locs
  end
end
