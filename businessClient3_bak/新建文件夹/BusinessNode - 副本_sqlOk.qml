import QtQuick 2.4
import QtQuick.LocalStorage 2.0

Rectangle {
    width: 360
    height: 360
    id: screen
    Text {
        id: textDisplay
        anchors.centerIn: parent
    }
    Component.onCompleted: {
        // 初始化数据库
        initialize();
        // 赋值
        setSetting("mySetting","ZX2");
        //获取一个值，并把它写在textDisplay里
        textDisplay.text = "The value of mySetting is:\n" + getSetting("mySetting");
    }


    function getDatabase() {
         return LocalStorage.openDatabaseSync("MyAppName", "1.0", "StorageDatabase", 100000);
    }

    // 程序打开时，初始化表
    function initialize() {
        var db = getDatabase();
        db.transaction(
            function(tx) {
                // 如果setting表不存在，则创建一个
                // 如果表存在，则跳过此步
                tx.executeSql('CREATE TABLE IF NOT EXISTS settings(setting TEXT UNIQUE, value TEXT)');
          });
    }

    // 插入数据
    function setSetting(setting, value) {
       var db = getDatabase();
       var res = "";
       db.transaction(function(tx) {
            var rs = tx.executeSql('INSERT OR REPLACE INTO settings VALUES (?,?);', [setting,value]);
                  //console.log(rs.rowsAffected)
                  if (rs.rowsAffected > 0) {
                    res = "OK";
                  } else {
                    res = "Error";
                  }
            }
      );
      console.log(res)
      return res;
    }

     // 获取数据
    function getSetting(setting) {
       var db = getDatabase();
       var res="";
       db.transaction(function(tx) {
         var rs = tx.executeSql('SELECT value FROM settings WHERE setting=?;', [setting]);
         if (rs.rows.length > 0) {
              res = rs.rows.item(0).value;
         } else {
             res = "Unknown";
         }
      })
      return res
    }
}
