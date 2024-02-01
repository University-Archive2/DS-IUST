package Tokenize;

import java.io.IOException;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.net.Socket;

class ClientHandler extends Thread {
    private final Socket socket;
    private final ObjectInputStream inputStream;
    private final ObjectOutputStream outputStream;

    ClientHandler(Socket socket, ObjectInputStream inputStream, ObjectOutputStream outputStream) {
        this.socket = socket;
        this.inputStream = inputStream;
        this.outputStream = outputStream;
    }

    @Override
    public void run() {
        while (true) {
            try {
                int index = inputStream.readInt();
                outputStream.writeObject(Constants.INPUT.substring(index * 4, index * 4 + 4));
                outputStream.flush();
            } catch (IOException e) {
                break;
            }
        }
    }
}
