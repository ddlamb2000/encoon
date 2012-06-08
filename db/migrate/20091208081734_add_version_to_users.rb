class AddVersionToUsers < ActiveRecord::Migration
  def self.up
    add_column :users, :version, :integer
    add_column :users, :begin, :date
  end

  def self.down
    remove_column :users, :version
    remove_column :users, :begin
  end
end
