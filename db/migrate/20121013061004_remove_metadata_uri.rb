class RemoveMetadataUri < ActiveRecord::Migration
  def up
    execute "delete from column_locs where uuid in (select column_uuid from column_mappings where db_column = 'uri')"
    execute "delete from columns where uuid in (select column_uuid from column_mappings where db_column = 'uri')"
    execute "delete from column_mappings where db_column = 'uri'"
  end

  def down
  end
end
