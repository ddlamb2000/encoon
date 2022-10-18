class AddActiveStatusToTables < ActiveRecord::Migration
  def self.up
    add_column :users, :active, :boolean
    add_column :workspaces, :active, :boolean
    add_column :grids, :active, :boolean
    add_column :columns, :active, :boolean
    add_column :rows, :active, :boolean

    User.find(:all).each do |item|
      item.active = true
      item.save! if item.valid?
    end
    Workspace.find(:all).each do |item|
      item.active = true
      item.save! if item.valid?
    end
    Grid.find(:all).each do |item|
      item.active = true
      item.save! if item.valid?
    end
    Column.find(:all).each do |item|
      item.active = true
      item.save! if item.valid?
    end
    Row.find(:all).each do |item|
      item.active = true
      item.save! if item.valid?
    end
  end

  def self.down
    remove_column :users, :active
    remove_column :workspaces, :active
    remove_column :grids, :active
    remove_column :columns, :active
    remove_column :rows, :active
  end
end
