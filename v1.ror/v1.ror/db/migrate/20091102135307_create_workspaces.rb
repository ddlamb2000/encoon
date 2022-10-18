class CreateWorkspaces < ActiveRecord::Migration
  def self.up
    create_table :workspaces do |t|
      t.date :begin
      t.date :end
      t.string :name
      t.text :description
      t.integer :version
      t.integer :revision

      t.timestamps
    end
  end

  def self.down
    drop_table :workspaces
  end
end
