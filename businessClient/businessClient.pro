#-------------------------------------------------
#
# Project created by QtCreator 2019-01-05T12:26:07
#
#-------------------------------------------------

QT       += core gui network websockets

greaterThan(QT_MAJOR_VERSION, 4): QT += widgets

TARGET = businessClient
TEMPLATE = app

# The following define makes your compiler emit warnings if you use
# any feature of Qt which has been marked as deprecated (the exact warnings
# depend on your compiler). Please consult the documentation of the
# deprecated API in order to know how to port your code away from it.
DEFINES += QT_DEPRECATED_WARNINGS

# You can also make your code fail to compile if you use deprecated APIs.
# In order to do so, uncomment the following line.
# You can also select to disable deprecated APIs only up to a certain version of Qt.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0


SOURCES += \
        main.cpp \
        mainwindow.cpp \
        mywebsock.cpp \
        protobuf/txdata.pb.cc \
        logindialog.cpp \
        dataexchanger.cpp

HEADERS += \
        mainwindow.h \
        mywebsock.h \
        protobuf/m2b.h \
        protobuf/txdata.pb.h \
        logindialog.h \
        dataexchanger.h

FORMS += \
        mainwindow.ui \
        logindialog.ui

CONFIG += mobility
MOBILITY = 

# 禁用(warning: unused parameter '变量名' [-Wunused-parameter])
QMAKE_CXXFLAGS += -Wno-unused-parameter

INCLUDEPATH += $$PWD/protobuf
INCLUDEPATH += $$PWD/protobuf/protobuf-3.6.1/src
# LIBS      += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.a
win32 {
# 让GUI程序(Run in terminal).
CONFIG      += console
LIBS        += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.mingw.a
}
android {
LIBS        += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.ndk.a
}
unix {
# LIBS      += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.unix.a
}
