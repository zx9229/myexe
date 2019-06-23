#include "mytts.h"
#include "myandroidcls.h"

MyTTS MyTTS::m_mytts;

MyTTS::MyTTS(QObject *parent) : QObject(parent), m_tts(nullptr)
{
    initialize();
}

MyTTS::~MyTTS()
{
}

void MyTTS::staticSpeak(const QString& text)
{
    m_mytts.speak(text);
}

void MyTTS::initialize()
{
#ifdef Q_OS_ANDROID
    android_tool::ttsInit();
#else
    m_tts = new QTextToSpeech(this);
#endif
}

void MyTTS::speak(const QString& text)
{
#ifdef Q_OS_ANDROID
    android_tool::ttsSpeak(text);
#else
    m_tts->say(text);
#endif
}
