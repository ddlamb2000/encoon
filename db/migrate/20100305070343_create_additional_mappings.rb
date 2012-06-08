require 'uuid'

class CreateAdditionalMappings < ActiveRecord::Migration
  def self.up
    uuid_gen = UUID.new

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => Column::ROOT_REFERENCE_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'grid_reference_uuid')

    GridMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :grid_uuid => GridMapping::ROOT_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_table => 'grid_mappings')
    
    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => GridMapping::ROOT_GRID_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'grid_uuid')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => GridMapping::ROOT_DB_TABLE,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'db_table')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => GridMapping::ROOT_DB_LOC_TABLE,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'db_loc_table')

    GridMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :grid_uuid => ColumnMapping::ROOT_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_table => 'column_mappings')
    
    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => ColumnMapping::ROOT_COLUMN_UUID,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'column_uuid')

    ColumnMapping.create!(
      :uuid => uuid_gen.generate,
      :begin => Entity.begin_of_time,
      :end => Entity.end_of_time,
      :version => 1,
      :enabled => true,
      :column_uuid => ColumnMapping::ROOT_DB_COLUMN,
      :create_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :update_user_uuid => 'eeba1320-dd45-012c-aafe-0026b0d63708',
      :db_column => 'db_column')

  end

  def self.down
  end
end
