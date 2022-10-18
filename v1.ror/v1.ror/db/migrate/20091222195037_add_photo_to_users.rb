class AddPhotoToUsers < ActiveRecord::Migration
  def self.up
    add_column :users, :photo, :binary, :limit => 1.megabyte
  end

  def self.down
    remove_column :users, :photo
  end
end