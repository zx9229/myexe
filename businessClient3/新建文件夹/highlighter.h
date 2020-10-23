#ifndef HIGHLIGHTER_H
#define HIGHLIGHTER_H
//https://doc.qt.io/qt-5/qtwidgets-richtext-syntaxhighlighter-example.html
#include <QSyntaxHighlighter>
#include <QRegularExpression>
#include <QQuickTextDocument>

class Highlighter : public QSyntaxHighlighter
{
    Q_OBJECT
    Q_PROPERTY(QTextDocument* document READ document WRITE setDocument)

public:
    //Highlighter(QTextDocument *parent = nullptr);
    Highlighter(QQuickTextDocument *parent = nullptr);

public:
    static void doQmlRegisterType();

protected:
    void highlightBlock(const QString &text) override;

private:
    struct HighlightingRule
    {
        QRegularExpression pattern;
        QTextCharFormat format;
    };
    QVector<HighlightingRule> highlightingRules;

    QRegularExpression commentStartExpression;
    QRegularExpression commentEndExpression;

    QTextCharFormat keywordFormat;
    QTextCharFormat classFormat;
    QTextCharFormat singleLineCommentFormat;
    QTextCharFormat multiLineCommentFormat;
    QTextCharFormat quotationFormat;
    QTextCharFormat functionFormat;
};

#endif // HIGHLIGHTER_H
