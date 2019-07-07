package zx.qtproject.example;

import android.content.Context;
import android.media.AudioManager;
import android.speech.tts.TextToSpeech;
import android.util.Log;
import java.text.MessageFormat;
import java.util.Locale;
//Android TTS 支持中文
//https://blog.csdn.net/u012532263/article/details/85338969

public class ZxTTS implements TextToSpeech.OnInitListener
{
    private static ZxTTS TTS = null;
    private static final String TAG = "ZxTTS";
    private Context      m_ctx = null;
    private TextToSpeech m_tts = null;
    public ZxTTS(Context ctx) {
        m_ctx = ctx;
        m_tts = new TextToSpeech(m_ctx, this);
    }
    @Override
    public void onInit(int status) {
        Log.v(TAG, MessageFormat.format("onInit, status={0}, SUCCESS={1}", TTS, TextToSpeech.SUCCESS));
        if (status == TextToSpeech.SUCCESS) {}
        if (false) {//它是将系统级别的媒体音量调至最大.
            AudioManager am = (AudioManager)m_ctx.getSystemService(Context.AUDIO_SERVICE);
            int amMusicStreamMaxVolume = am.getStreamMaxVolume(AudioManager.STREAM_MUSIC);
            am.setStreamVolume(AudioManager.STREAM_MUSIC, amMusicStreamMaxVolume, 0);
        }
    }
    public boolean speak(final String text) {
        //public int speak (String text, int queueMode, HashMap<String, String> params) // Deprecated in API level 21
        //public int speak (CharSequence text, int queueMode, Bundle params, String utteranceId)
        int retVal = m_tts.speak(text, TextToSpeech.QUEUE_FLUSH, null);
        return (retVal == TextToSpeech.SUCCESS);
    }
    public static void staticInit(Context ctx) {
        if (TTS == null) {
            TTS = new ZxTTS(ctx);
        }
        Log.v(TAG, MessageFormat.format("staticInit, TTS={0}", TTS));
    }
    public static boolean staticSpeak(final String text) {
        boolean retVal = (TTS == null) ? false : TTS.speak(text);
        //Log.v(TAG, MessageFormat.format("staticSpeak, {0}, TTS={1}, text={2}", retVal, TTS, text));
        return retVal;
    }
}
