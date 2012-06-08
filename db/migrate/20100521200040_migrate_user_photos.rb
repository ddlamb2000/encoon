class MigrateUserPhotos < ActiveRecord::Migration
  def self.up
    execute "insert into row_attachments(uuid, document, content_type) select uuid, photo, 'image/?' from users where photo is not null"
  end

  def self.down
  end
end
