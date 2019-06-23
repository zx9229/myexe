#ifndef MY_TTS_H_H
#define MY_TTS_H_H
#include <QObject>
#include <QTextToSpeech>

class MyTTS : public QObject
{
    Q_OBJECT
public:
    ~MyTTS();
    static void staticSpeak(const QString& text);
private:
    explicit MyTTS(QObject *parent = nullptr);
    void initialize();
    void speak(const QString& text);
private:
    QTextToSpeech* m_tts;
    static MyTTS   m_mytts;
};

#endif//MY_TTS_H_H
