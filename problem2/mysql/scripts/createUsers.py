# coding:utf-8
from faker import Factory

# 苗字と名前のリストから名前を生成
def randomName():
    f = Factory.create('ja_JP')
    return f.name()

# 出力するファイル名
OUTPUT_FILE = "../UsersTestData.sql"

# 登録するデータ件数
RECORD_COUNT = 30

# 実行するSQLコマンド文字列
sqlCommands = ""

# 使用するデータベースを指定(今回はCreateTestData)
sqlCommands += "USE app;\n"

# 登録するデータの数だけINSERT文を生成
for i in range(RECORD_COUNT):

    # 登録するランダムなデータの生成
    id  = i
    name = randomName()

    # ランダムなデータからInsert文を生成
    sqlCommands += "INSERT INTO users " \
                   "(id, user_id, name) " \
                   "VALUES (0, '{}', '{}');\n"\
                   .format(id, name)

# 生成したSQLコマンドをファイルに書き出す
f = open(OUTPUT_FILE, 'w')
f.write(sqlCommands)
f.close()
