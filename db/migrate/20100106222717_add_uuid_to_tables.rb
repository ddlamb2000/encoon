require 'uuid'

class AddUuidToTables < ActiveRecord::Migration
  def self.up
    add_column :users, :uuid, :string, :limit => 36
    add_column :workspaces, :uuid, :string, :limit => 36
    add_column :grids, :uuid, :string, :limit => 36
    add_column :columns, :uuid, :string, :limit => 36
    add_column :rows, :uuid, :string, :limit => 36
    
    uuid_gen = UUID.new

    User.find(:all).each do |item|
      if item.uuid.blank?
        item.uuid = uuid_gen.generate
        item.save! if item.valid?
      end
    end
    Workspace.find(:all).each do |item|
      if item.uuid.blank?
        item.uuid = uuid_gen.generate
        item.save! if item.valid?
      end
    end
    Grid.find(:all).each do |item|
      if item.uuid.blank?
        item.uuid = uuid_gen.generate
        item.save! if item.valid?
      end
    end
    Column.find(:all).each do |item|
      if item.uuid.blank?
        item.uuid = uuid_gen.generate
        item.save! if item.valid?
      end
    end
    Row.find(:all).each do |item|
      if item.uuid.blank?
        item.uuid = uuid_gen.generate
        item.save! if item.valid?
      end
    end
  end

  def self.down
    remove_column :users, :uuid
    remove_column :workspaces, :uuid
    remove_column :grids, :uuid
    remove_column :columns, :uuid
    remove_column :rows, :uuid
  end
end