class CreateWorkspacesLoc < ActiveRecord::Migration
  def self.up
    create_table :workspace_locs do |t|
      t.string   :uuid, :limit => 36
      t.integer  :version
      t.string   :locale, :limit => 10
      t.string   :base_locale, :limit => 10
      t.string   :name
      t.text     :description
      t.integer  :lock_version, :default => 0
    end

    add_index :workspace_locs, [:uuid, :version, :locale]

    Workspace.find(:all).each do |item|
      puts "Update workspace " + item.id.to_s
      LANGUAGES.each do |lang, locale|
        puts "Update locale " + locale.to_s
        workspace_loc = WorkspaceLoc.new
        workspace_loc.uuid = item.uuid
        workspace_loc.version = item.version
        workspace_loc.locale = locale
        workspace_loc.base_locale = I18n.default_locale.to_s
        workspace_loc.name = item.name + (locale != "en" ? " (" + locale + ")" : "")
        workspace_loc.description = item.description + (locale != "en" ? " (" + locale + ")" : "") if item.description.present?
        workspace_loc.save!
      end
    end
end

  def self.down
    drop_table :workspace_locs
  end
end