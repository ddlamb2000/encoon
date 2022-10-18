class MigrateUserPasswords < ActiveRecord::Migration
  def self.up
    execute "insert into row_passwords(uuid, salt, password) select uuid, salt, password from users where password is not null"
  end

  def self.down
  end
end
