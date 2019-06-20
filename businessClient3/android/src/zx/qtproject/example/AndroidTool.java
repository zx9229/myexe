//因为【package="zx.qtproject.example"】所以:
//java file goes in android/src/zx/qtproject/example/AndroidTool.java
package zx.qtproject.example;

import android.content.Context;
import android.os.Looper;
import android.os.Process;//Log.v("TAG", MessageFormat.format("myPid={0}", Process.myPid()))=
import android.util.Log;
import android.widget.Toast;
import java.text.MessageFormat;//Log.v("TAG", MessageFormat.format("myPid={0}", Process.myPid()))=

public class AndroidTool
{
    public static void logVerbose(String tag, String message) {
        Log.v(tag, message);
    }

    public static void toastShow(Context context, String message) {
        //Log.v("My_AndroidTool", MessageFormat.format("toastShow, context={0}, message={1}", context, message));
        try {
            Toast.makeText(context, message, Toast.LENGTH_LONG).show();
        }catch(Exception e) {
            //Looper.prepare();
            //Toast.makeText(context, message, Toast.LENGTH_LONG).show();
            //Looper.myLooper().quitSafely();
            //Looper.loop();
            //Log.v("My_AndroidTool", MessageFormat.format("toastShow, e={0}, will return", e));
        }
    }
}
//Object.toString()
