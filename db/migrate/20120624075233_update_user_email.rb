class UpdateUserEmail < ActiveRecord::Migration
  def up
    add_index :users, :email
    execute "update users set email = identifier"
  end

  def down
  end
end