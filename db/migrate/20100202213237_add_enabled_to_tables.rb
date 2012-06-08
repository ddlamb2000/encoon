class AddEnabledToTables < ActiveRecord::Migration
  def self.up
    add_column :users, :enabled, :boolean
    add_column :workspaces, :enabled, :boolean
    add_column :grids, :enabled, :boolean
    add_column :columns, :enabled, :boolean
    add_column :rows, :enabled, :boolean

    User.find(:all).each do |item|
      item.enabled = true
      item.save! if item.valid?
    end
    Workspace.find(:all).each do |item|
      item.enabled = true
      item.save! if item.valid?
    end
    Grid.find(:all).each do |item|
      item.enabled = true
      item.save! if item.valid?
    end
    Column.find(:all).each do |item|
      item.enabled = true
      item.save! if item.valid?
    end
    Row.find(:all).each do |item|
      item.enabled = true
      item.save! if item.valid?
    end

    remove_column :users, :active
    remove_column :workspaces, :active
    remove_column :grids, :active
    remove_column :columns, :active
    remove_column :rows, :active
  end

  def self.down
    remove_column :users, :enabled
    remove_column :workspaces, :enabled
    remove_column :grids, :enabled
    remove_column :columns, :enabled
    remove_column :rows, :enabled

    add_column :users, :active, :boolean
    add_column :workspaces, :active, :boolean
    add_column :grids, :active, :boolean
    add_column :columns, :active, :boolean
    add_column :rows, :active, :boolean
  end
end
