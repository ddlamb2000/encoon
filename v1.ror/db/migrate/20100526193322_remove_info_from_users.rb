class RemoveInfoFromUsers < ActiveRecord::Migration
  def self.up
    remove_column :users, :salt
    remove_column :users, :password
    remove_column :users, :photo
  end

  def self.down
    add_column :users, :salt, :string
    add_column :users, :password, :string
    add_column :users, :photo, :binary, :limit => 1.megabyte
  end
end
