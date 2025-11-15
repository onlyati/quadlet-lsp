package commonlog

import (
	"bufio"
	"net"

	"github.com/tliron/kutil/util"
)

//
// LoggerServer
//

// A TCP server file that forwards all lines written to it to a [Logger].
type LoggerServer struct {
	IPStack util.IPStack
	Address string
	Port    int
	Log     Logger
	Level   Level

	ClientAddressPorts []string

	listeners []net.Listener
}

func NewLoggerServer(ipStack util.IPStack, address string, port int, log Logger, level Level) *LoggerServer {
	return &LoggerServer{
		IPStack: ipStack,
		Address: address,
		Port:    port,
		Log:     NewKeyValueLogger(log, "tcp", port),
		Level:   level,
	}
}

func (self *LoggerServer) Start() error {
	return self.IPStack.StartServers(self.Address, self.start)
}

func (self *LoggerServer) Stop() {
	for index, listener := range self.listeners {
		self.Log.Notice("stopping logger server",
			"index", index)
		if err := listener.Close(); err != nil {
			self.Log.Error(err.Error())
		}
		self.Log.Notice("stopped logger server",
			"index", index)
	}
}

// ([util.IPStackStartServerFunc] signature)
func (self *LoggerServer) start(level2protocol string, address string) error {
	if address, err := util.ToReachableIPAddress(address); err == nil {
		addressPort := util.JoinIPAddressPort(address, self.Port)
		if listener, err := net.Listen(level2protocol, addressPort); err == nil {
			self.Log.Notice("starting logger server",
				"index", len(self.listeners),
				"level2protocol", level2protocol,
				"addressPort", listener.Addr().String())

			self.ClientAddressPorts = append(self.ClientAddressPorts, util.IPAddressPortWithoutZone(addressPort))
			self.listeners = append(self.listeners, listener)

			go func() {
				for {
					if conn, err := listener.Accept(); err == nil {
						self.Log.Debug("accepted logger server connection")
						go self.handle(conn)
					} else {
						self.Log.Critical(err.Error())
						return
					}
				}
			}()

			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *LoggerServer) handle(conn net.Conn) {
	defer CallAndLogError(func() error {
		self.Log.Debug("closing logger server connection")
		return conn.Close()
	}, "Conn.Close", self.Log)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		self.Log.Log(self.Level, 0, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		self.Log.Error(err.Error())
	}
}
