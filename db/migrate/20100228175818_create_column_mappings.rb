require 'uuid'

class CreateColumnMappings < ActiveRecord::Migration
  def self.up
    create_table :column_mappings, :force => true do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.integer  :lock_version,        :default => 0
      t.string   :column_uuid,           :limit => 36
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.string   :db_column
      t.timestamps
    end

    add_index :column_mappings, [:column_uuid]
  end

  def self.down
    drop_table :column_mappings
  end
end
