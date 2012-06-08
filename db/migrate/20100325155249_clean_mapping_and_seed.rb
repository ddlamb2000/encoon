class CleanMappingAndSeed < ActiveRecord::Migration
  def self.up
    GridMapping.destroy_all
    ColumnMapping.destroy_all
    upload = Upload.create
    upload.data_file=File.new("db/system.xml")
  end

  def self.down
  end
end
