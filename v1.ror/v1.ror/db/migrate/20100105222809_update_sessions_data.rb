class UpdateSessionsData < ActiveRecord::Migration
  def self.up
    change_column :sessions, :data, :text, :limit => 1.megabyte
  end

  def self.down
    change_column :sessions, :data, :text
  end
end
