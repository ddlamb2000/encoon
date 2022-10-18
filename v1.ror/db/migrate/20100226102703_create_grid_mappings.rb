require 'uuid'

class CreateGridMappings < ActiveRecord::Migration
  def self.up
    create_table :grid_mappings, :force => true do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.integer  :lock_version,        :default => 0
      t.string   :grid_uuid,           :limit => 36
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.string   :db_table
      t.string   :db_loc_table
      t.timestamps
    end

    add_index :grid_mappings, [:grid_uuid], :name => "index_grid_mappings_on_grid_uuid"
  end

  def self.down
    drop_table :grid_mappings
  end
end