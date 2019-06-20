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

    public static void funTest(Context ctx) {
        Log.v(TAG, MessageFormat.format("myPid={0}, funTest, beg", Process.myPid()));
        ctx.startForegroundService(new Intent(ctx, ZxService.class));
        Log.v(TAG, MessageFormat.format("myPid={0}, funTest, end", Process.myPid()));
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        Log.v(TAG, MessageFormat.format("myPid={0}, onCreate, savedInstanceState={0}, beg", Process.myPid(), savedInstanceState));
        super.onCreate(savedInstanceState);
        if (false) {
            Intent intent = new Intent(this, ZxService.class);
            startService(intent);
        }
        Log.v(TAG, MessageFormat.format("myPid={0}, onCreate, savedInstanceState={0}, end", Process.myPid(), savedInstanceState));
    }
}
