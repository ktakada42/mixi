# coding:utf-8
import random

# 出力するファイル名
OUTPUT_FILE = "../BlockListTestData.sql"

# 登録するデータ件数
RECORD_COUNT = 30

# 実行するSQLコマンド文字列
sqlCommands = ""

# 使用するデータベースを指定(今回はCreateTestData)
sqlCommands += "USE app;\n"

# 登録するデータの数だけINSERT文を生成
for i in range(2, RECORD_COUNT):
    ns = []
    cnt = 0
    regi = random.randint(0, 4)

    # 登録するランダムなデータの生成
    id1 = i
    # 0~100の乱数を生成
    while cnt < regi:
        id2 = random.randint(0, 29)
        if not id2 in ns and id1 != id2:
            ns.append(id2)
            # ランダムなデータからInsert文を生成
            sqlCommands += "INSERT INTO block_list " \
                           "(id, user1_id, user2_id) " \
                           "VALUES (0, '{}', '{}');\n" \
                .format(id1, id2)
        cnt += 1

# 生成したSQLコマンドをファイルに書き出す
f = open(OUTPUT_FILE, 'w')
f.write(sqlCommands)
f.close()
