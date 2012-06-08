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
    
    uuid_gen = UUID.new

    grid_mapping = GridMapping.create!(
                      :uuid => uuid_gen.generate,
                      :begin => Entity.begin_of_time,
                      :end => Entity.end_of_time,
                      :version => 1,
                      :enabled => true,
                      :grid_uuid => Workspace::ROOT_UUID,
                      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :db_table => 'workspaces',
                      :db_loc_table => 'workspace_locs')
    
    grid_mapping = GridMapping.create!(
                      :uuid => uuid_gen.generate,
                      :begin => Entity.begin_of_time,
                      :end => Entity.end_of_time,
                      :version => 1,
                      :enabled => true,
                      :grid_uuid => Grid::ROOT_UUID,
                      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :db_table => 'grids',
                      :db_loc_table => 'grid_locs')

    grid_mapping = GridMapping.create!(
                      :uuid => uuid_gen.generate,
                      :begin => Entity.begin_of_time,
                      :end => Entity.end_of_time,
                      :version => 1,
                      :enabled => true,
                      :grid_uuid => Column::ROOT_UUID,
                      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :db_table => 'columns',
                      :db_loc_table => 'column_locs')

    grid_mapping = GridMapping.create!(
                      :uuid => uuid_gen.generate,
                      :begin => Entity.begin_of_time,
                      :end => Entity.end_of_time,
                      :version => 1,
                      :enabled => true,
                      :grid_uuid => User::ROOT_UUID,
                      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
                      :db_table => 'users')
  end

  def self.down
    drop_table :grid_mappings
  end
end