class CreateGrids < ActiveRecord::Migration
  def self.up
    create_table :grids do |t|
      t.date :begin
      t.date :end
      t.string :name
      t.text :description
      t.integer :version
      t.integer :revision
      t.references :workspace
      t.string :kind

      t.timestamps
    end
  end

  def self.down
    drop_table :grids
  end
end
