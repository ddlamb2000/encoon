class AddEndToUsers < ActiveRecord::Migration
  def self.up
    add_column :users, :end, :date
  end

  def self.down
    remove_column :users, :end
  end
end
