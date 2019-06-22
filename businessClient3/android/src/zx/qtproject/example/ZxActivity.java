package zx.qtproject.example;

import zx.qtproject.example.ZxService;

import org.qtproject.qt5.android.bindings.QtActivity;

import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.ServiceConnection;
import android.os.Bundle;
import android.os.IBinder;
import android.os.Process;//Log.v("TAG", MessageFormat.format("myPid={0}", Process.myPid()))=
import android.util.Log;
import android.widget.Toast;
import java.text.MessageFormat;//Log.v("TAG", MessageFormat.format("myPid={0}", Process.myPid()))=

public class ZxActivity extends QtActivity
{
    private static final String TAG = "ZxActivity";

    public static void startTheService(Context ctx) {
        //Log.v(TAG, "startTheService, beg");
        ctx.startForegroundService(new Intent(ctx, ZxService.class));
        //Log.v(TAG, "startTheService, end");
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        Log.v(TAG, MessageFormat.format("myPid={0}, onCreate, savedInstanceState={0}, beg", Process.myPid(), savedInstanceState));
        super.onCreate(savedInstanceState);
        if (false) { // 这里仅用来表达:可以在启动程序的时候,立即启动服务.
            Intent intent = new Intent(this, ZxService.class);
            startService(intent);
        }
        Log.v(TAG, MessageFormat.format("myPid={0}, onCreate, savedInstanceState={0}, end", Process.myPid(), savedInstanceState));
    }
}
