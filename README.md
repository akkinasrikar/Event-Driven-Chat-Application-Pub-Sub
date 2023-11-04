# Building an Event-Driven Chat Application (Similar to Whatsapp) with Go Channels, Mutex, and Goroutines

Description:

In the world of real-time communication, designing a robust chat application that can handle concurrent user interactions is a challenging task. Event-Driven Architecture (EDA) is an approach that can greatly enhance the scalability and responsiveness of such applications. This description explores the creation of an EDA-based chat application using the Go programming language, leveraging Go channels, mutexes, and goroutines to build a secure and efficient messaging platform.

**Key Components:**

1. **Pub-Sub Model:** Our chat application adopts a publish-subscribe model where messages are published to channels and subscribers receive them in real-time. This model allows for seamless communication between users without the need for polling.

2. **Go Channels:** Go channels are at the heart of our application, facilitating communication between different parts of the system. We use channels to send and receive messages, manage user connections, and handle events asynchronously.

3. **Mutexes:** To ensure thread safety and prevent race conditions, we employ mutexes to control access to shared resources like the user list and chat history. This guarantees data consistency and avoids conflicts when multiple users interact with the application simultaneously.

4. **Goroutines:** Goroutines are lightweight threads in Go that enable concurrent execution. We use goroutines to manage user connections, handle incoming messages, and broadcast messages to all connected users. This concurrent processing ensures optimal performance even with a large number of users.

The "Event Driven Messenger" project is a real-time communication application designed with a focus on leveraging Go routines and channels, following event-driven pub-sub system design principles. Unlike traditional messenger applications, this project is volatile in nature, meaning it does not rely on a persistent database, enhancing its performance and speed.

**Key Features:**

1. **Multiple Admins:** The system allows for multiple administrators to manage and moderate conversations and groups.

2. **Leave Group:** Users can leave a group at any time, providing flexibility in managing their participation.

3. **Unsubscribe User:** The ability to unsubscribe a user from specific channels or groups, ensuring user preferences are respected.

4. **Group Limit:** Admins can set limits on the number of participants in a group, maintaining the group's intended size.

5. **Group Admin:** Assigning group administrators with specific privileges for managing group activities.

6. **History of the Conversation:** The application stores and provides access to the history of conversations for users to refer back to previous discussions.

7. **One-to-One Messaging:** Users can have private conversations with individuals, ensuring secure and direct communication.

8. **One-to-Many Broadcasting:** Broadcasting messages to multiple users or groups with a single action, facilitating mass communication.

9. **Message Broker:** The system includes a message broker that manages message distribution and ensures efficient communication between users and groups.

10. **Publisher-Subscriber Model:** Utilizes a publisher-subscriber model for event-driven communication, allowing users to subscribe to specific topics or channels and receive real-time updates.

This project represents a dynamic and efficient event-driven messenger application, eliminating the need for a database while focusing on Go routines and channels to deliver fast and responsive communication. It caters to the needs of various users and administrators, offering a wide range of features for effective communication and collaboration.
