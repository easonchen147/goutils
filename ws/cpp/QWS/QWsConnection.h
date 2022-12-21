#pragma once
#include <QtCore>
#include <QWebSocket>

typedef unsigned char BYTE;

class QWsConnection : public QObject
{
    Q_OBJECT

    using MsgHandler = std::function<void(QWebSocket*, const QByteArray&)>;
    using EvtHandler = std::function<void(QWebSocket*)>;
    using ErrHandler = std::function<void(QWebSocket*, QAbstractSocket::SocketError)>;
    using DisplacedHandler = std::function<void(QWebSocket*, QString, QString, int64_t)>;

public:
    explicit QWsConnection(QObject* parent = nullptr);
    ~QWsConnection();

    void RegisterMsgHandler(uint32_t protocolId, MsgHandler handler);

    inline void SetEstablishHandler(EvtHandler establishHandler) { m_establishHandler = establishHandler; }
    inline void SetCloseHandler(EvtHandler closeHandler) { m_closeHandler = closeHandler; }
    inline void SetErrHandler(ErrHandler errHandler) { m_errHandler = errHandler; }
    inline void SetDisplacedHandler(DisplacedHandler displacedHandler) { m_displacedHandler = displacedHandler; }

    void Connect(const QString& url, int retryInterval = 0);
    void SendMsg(uint32_t protocolId, const QByteArray& data);

protected:
    struct innerMsgPack
    {
        BYTE        packetHeadFlag[2];
        uint32_t    packetLength;
        uint32_t    protocolId;
        QByteArray  dataBuffer;

        innerMsgPack():packetLength(0), protocolId(0)
        {
            packetHeadFlag[2] = { 0 };
        }
    };

    QByteArray _PackMsg(uint32_t protocolId, QByteArray dataBuffer);
    innerMsgPack _UnpackMsg(QByteArray rawMsg);
    void _OnDisplaced(QWebSocket* ws, QByteArray msgData);

private:
    QWebSocket*                 m_pWs;
    bool                        m_bConnected;

    QHash<uint32_t, MsgHandler> m_mapMsgHandler;
    EvtHandler                  m_establishHandler;
    EvtHandler                  m_closeHandler;
    ErrHandler                  m_errHandler;
    DisplacedHandler            m_displacedHandler;

    static BYTE                 m_packetHeadFlag[2];
};
