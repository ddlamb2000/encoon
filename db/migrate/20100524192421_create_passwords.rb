class CreatePasswords < ActiveRecord::Migration
  def self.up
    create_table :row_passwords do |t|
      t.string   :uuid, :limit => 36
      t.string   :password
      t.string   :salt
      t.integer  :lock_version, :default => 0
    end

    add_index :row_passwords, [:uuid]
  end

  def self.down
    drop_table :row_passwords
  end
end
