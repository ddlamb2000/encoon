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
    
    uuid_gen = UUID.new

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Grid::ROOT_WORKSPACE_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'workspace_uuid')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Column::ROOT_GRID_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'grid_uuid')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Column::ROOT_KIND_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'kind')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Column::ROOT_NUMBER_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'number')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Column::ROOT_DISPLAY_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'display')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Column::ROOT_TRANSLATED_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'translated')
  end

  def self.down
    drop_table :column_mappings
  end
end
