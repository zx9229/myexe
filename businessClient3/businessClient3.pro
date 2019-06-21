QT += quick core network websockets sql remoteobjects
android {
QT += androidextras
}
CONFIG += c++11

REPC_SOURCE  = dataro.rep
REPC_REPLICA = dataro.rep
# The following define makes your compiler emit warnings if you use
# any feature of Qt which as been marked deprecated (the exact warnings
# depend on your compiler). Please consult the documentation of the
# deprecated API in order to know how to port your code away from it.
DEFINES += QT_DEPRECATED_WARNINGS

# You can also make your code fail to compile if you use deprecated APIs.
# In order to do so, uncomment the following line.
# You can also select to disable deprecated APIs only up to a certain version of Qt.
#DEFINES += QT_DISABLE_DEPRECATED_BEFORE=0x060000    # disables all the APIs deprecated before Qt 6.0.0

SOURCES += \
    dataexchanger.cpp \
    main.cpp \
    mywebsock.cpp \
    protobuf/txdata.pb.cc \
    zxtools.cpp \
    datarosvr.cpp \
    datawrapper.cpp \
    temputils.cpp \
    myandroidcls.cpp

RESOURCES += qml.qrc

# Additional import path used to resolve QML modules in Qt Creator's code model
QML_IMPORT_PATH =

# Additional import path used to resolve QML modules just for Qt Quick Designer
QML_DESIGNER_IMPORT_PATH =

# Default rules for deployment.
qnx: target.path = /tmp/$${TARGET}/bin
else: unix:!android: target.path = /opt/$${TARGET}/bin
!isEmpty(target.path): INSTALLS += target

HEADERS += \
    dataexchanger.h \
    mysqltablemodel.h \
    mywebsock.h \
    protobuf/m2b.h \
    protobuf/txdata.pb.h \
    safesynchcache.h \
    sqlstruct.h \
    zxtools.h \
    datarosvr.h \
    datawrapper.h \
    mylog.h \
    temputils.h \
    myandroidcls.h

# 禁用(warning: unused parameter '变量名' [-Wunused-parameter])
QMAKE_CXXFLAGS += -Wno-unused-parameter

INCLUDEPATH += $$PWD/protobuf
INCLUDEPATH += $$PWD/protobuf/protobuf-3.6.1/src
# LIBS      += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.a
win32 {
LIBS        += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.mingw.a
}
android {
LIBS        += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.ndk.a
}
unix {
# LIBS      += $$PWD/protobuf/protobuf-3.6.1/lib/libprotobuf.unix.a
}

DISTFILES += \
    android/AndroidManifest.xml \
    android/gradle/wrapper/gradle-wrapper.jar \
    android/gradlew \
    android/res/values/libs.xml \
    android/build.gradle \
    android/gradle/wrapper/gradle-wrapper.properties \
    android/gradlew.bat \
    dataro.rep \
    android/src/zx/qtproject/example/AndroidTool.java \
    android/src/zx/qtproject/example/ZxActivity.java \
    android/src/zx/qtproject/example/ZxService.java

ANDROID_PACKAGE_SOURCE_DIR = $$PWD/android
