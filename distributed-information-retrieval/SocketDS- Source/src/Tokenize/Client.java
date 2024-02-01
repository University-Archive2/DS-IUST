package Tokenize;

import java.io.*;
import java.net.*;
import java.util.Scanner;

public class Client {
    public static void main(String[] args) throws IOException, ClassNotFoundException {

        Scanner scanner = new Scanner(System.in);

        Socket s = new Socket("localhost", 3000);

        ObjectOutputStream outputStream = new ObjectOutputStream(s.getOutputStream());
        ObjectInputStream inputStream = new ObjectInputStream(s.getInputStream());

        int numOfPackets = Constants.INPUT.length() / 4;
        String[] output = new String[numOfPackets];

        for (int i = 0; i < numOfPackets; i++) {
            int index = Integer.parseInt(scanner.next());
            outputStream.writeInt(index);
            outputStream.flush();
            String response = (String) inputStream.readObject();
            output[index] = response;
            System.out.println(response);
        }
        System.out.println(String.join("", output));
    }
}
