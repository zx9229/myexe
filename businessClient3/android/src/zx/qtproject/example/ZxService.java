package zx.qtproject.example;

import org.qtproject.qt5.android.bindings.QtService;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.os.Build;
import android.os.IBinder;
import android.os.Process;//Log.v("TAG", MessageFormat.format("myPid={0}", Process.myPid()))=
import android.util.Log;
import android.widget.Toast;
import java.text.MessageFormat;//Log.v("TAG", MessageFormat.format("myPid={0}", Process.myPid()))=
import java.lang.reflect.Method;
import android.widget.RemoteViews;
import java.text.SimpleDateFormat;

// Service | Android Developers
// https://developer.android.com/reference/android/app/Service

public class ZxService extends QtService
{
    private static final String TAG = "ZxService";
    private static final int idINT = 1;

    public static void funTest(Context ctx) {
    }

    //首次创建服务时，系统将调用此方法来执行一次性设置程序（在调用 onStartCommand() 或 onBind() 之前）.
    //如果服务已在运行，则不会调用此方法。该方法只被调用一次.
    //Called by the system when the service is first created.
    @Override
    public void onCreate() {
        Log.v(TAG, "onCreate, beg");
        super.onCreate();
        startForeground(this.idINT, getNotification());
        Log.v(TAG, "onCreate, end");
    }

    //This method was deprecated in API level 15. Implement onStartCommand(android.content.Intent, int, int) instead.
    //public int onStart(Intent intent, int startId) { return super.onStart(intent, startId); }

    //每次通过startService方法启动Service时都会被回调.
    //Called by the system every time a client explicitly starts the service by calling Context.startService(Intent), providing the arguments it supplied and a unique integer token representing the start request.
    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        Log.v(TAG, MessageFormat.format("onStartCommand, intent={0}, flags={1}, startId={2}", intent, flags, startId));
        return super.onStartCommand(intent, flags, startId);
    }

    //当另一个组件想通过调用 bindService() 与服务绑定（例如执行 RPC）时，系统将调用此方法.
    //Return the communication channel to the service.
    @Override
    public IBinder onBind(Intent intent) {
        Log.v(TAG, MessageFormat.format("onBind, intent={0}", intent));
        return super.onBind(intent);
    }

    @Override
    public boolean onUnbind (Intent intent) {
        Log.v(TAG, MessageFormat.format("onUnbind, intent={0}", intent));
        return super.onUnbind(intent);
    }

    //Called by the system to notify a Service that it is no longer used and is being removed.
    @Override
    public void onDestroy() {
        super.onDestroy();
        Log.v(TAG, "onDestroy, called");
    }

    private Notification getNotification() {
        //https://developer.android.com/reference/android/app/Notification.Builder
        //https://developer.android.com/reference/android/R
        final String my_ChannelId = "MyService_CHANNEL_ID";
        final String my_CHANNEL_NAME = "MyService_CHANNEL_NAME";
        final String my_title =  MessageFormat.format("测试标题:[{0}]", new SimpleDateFormat("yyyy-MM-dd HH:mm:ss").format(System.currentTimeMillis()));
        final String my_text = "ZX测试文本";
        NotificationManager nManager = (NotificationManager)getSystemService(NOTIFICATION_SERVICE);
        NotificationChannel nChannel = new NotificationChannel(my_ChannelId, my_CHANNEL_NAME, NotificationManager.IMPORTANCE_HIGH);
        if (true) {
            nChannel.enableLights(true);//设置提示灯.
            nChannel.setLightColor(Color.RED);//设置提示灯颜色.
            nChannel.setShowBadge(true);//显示logo.
            nChannel.setDescription("My测试描述");//Sets the user visible description of this channel.
            nChannel.setLockscreenVisibility(Notification.VISIBILITY_PUBLIC); //设置锁屏可见 VISIBILITY_PUBLIC=可见.
            nChannel.enableLights(true); //是否在桌面icon右上角展示小红点.
            nChannel.setLightColor(Color.YELLOW); //小红点颜色.
            nChannel.setShowBadge(true); //是否在久按桌面图标时显示此渠道的通知.
        }
        nManager.createNotificationChannel(nChannel);
        Notification notification = new Notification.Builder(this)
            .setChannelId(my_ChannelId)
            .setContentTitle(my_title)//标题.
            //.setContentText(my_text)//内容.
            .setWhen(System.currentTimeMillis())
            .setShowWhen(true)
            .setSmallIcon(android.R.drawable.ic_dialog_info)//小图标一定需要设置,否则会报错(如果不设置它启动服务前台化不会报错,但是你会发现这个通知不会启动),如果是普通通知,不设置必然报错.
            //.setLargeIcon(android.graphics.drawable.Icon.createWithResource("android.R.drawable", android.R.drawable.ic_menu_help))
            .build();
        if(false) {
            RemoteViews ntfSmall = new RemoteViews(getPackageName(), android.R.drawable.ic_dialog_info);
            notification.contentView = ntfSmall;
        }
        return notification;
    }

    public static void testStartForegroundService(Context ctx) {
        ctx.startForegroundService(new Intent(ctx, ZxService.class));
    }

    public void testStartForeground() {
        this.startForeground(this.idINT, getNotification());
    }

    public void testStopForeground() {
        if(true) {
            this.stopForeground(this.idINT);
        } else {
            this.stopForeground(true);
        }
    }
}
