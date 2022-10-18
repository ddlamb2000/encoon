class AddGridUuidToAudits < ActiveRecord::Migration
  def self.up
    add_column :audits, :grid_uuid, :string, :limit => 36
    add_index :audits, [:update_user_uuid]
  end

  def self.down
    remove_column :audits, :grid_uuid
  end
end
