class RemoveGridUuidFromAudits < ActiveRecord::Migration
  def self.up
    remove_column :audits, :row_uuid
    remove_column :audits, :grid_uuid
    remove_column :audits, :begin
    remove_column :audits, :end
    remove_column :audits, :enabled
    remove_column :audits, :created_at
    remove_column :audits, :create_user_uuid
    add_column :audits, :locale, :string
    add_index :audits, [:uuid]
    remove_index "audits", :name => "index_audits_on_grid_uuid_and_row_uuid_and_version"
end

  def self.down
    add_column :audits, :row_uuid, :string
    add_column :audits, :grid_uuid, :string
    add_column :audits, :begin, :string
    add_column :audits, :end, :string
    add_column :audits, :enabled, :string
    add_column :audits, :created_at, :string
    add_column :audits, :create_user_uuid, :string
    remove_column :audits, :locale
    remove_index :audits, [:uuid]
    add_index "audits", ["version"], :name => "index_audits_on_grid_uuid_and_row_uuid_and_version"
  end
end