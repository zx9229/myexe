import QtQuick 2.4
import QtQuick.LocalStorage 2.0

BusinessNodeForm {
    id:bnf
    function loadImageData(){
        var db = openDatabaseSync("MyDB","1.0","my model SQL",5000)
        db.transaction(
            function(tx){
                tx.executeSql('CREATE TABLE IF NOT EXISTS Image(id INTEGER primary key, title TEXT,picture TEXT)');
                var rs = tx.executeSql('SELECT * FROM Image');
                if(rs.rows.length > 0){
                    var index = 0;
                    while(index<rs.rows.length){
                        var myItem = rs.rows.item(index);
                        bnf.mymodel.append({"id":myItem.id,"title":myItem.title,"picture":myItem.picture});
                        index++;
                    }
                }else{
                    bnf.mymodel.append({"id":1,"title":"apple","picture":"content/pics/apple.png"});
                    bnf.mymodel.append({"id":2,"title":"banne","picture":"content/pics/Qt.png"});
                }
            }
        )
    }
    function saveImageData(){
        var db = openDatabaseSync("MyDB","1.0","my model SQL",5000)
        db.transaction(
            function(tx){
                tx.executeSql('DROP TABLE Image');
                tx.executeSql('CREATE TABLE IF NOT EXISTS Image(id INTEGER primary key, title TEXT,picture TEXT)');
                var index = 0;
                while(index<bnf.mymodel.count){
                    var myItem = bnf.mymodel.get(index);
                    tx.executeSql('INSERT INTO Image VALUES(?,?,?)',[myItem.id,myItem.title,myItem.picture]);
                    index++;
                }
            }
        )
    }
    Connections{
        target: bnf.mymodel
        Component.onCompleted: loadImageData()
    }
}
