# Jepsen consistency model

Jepsen consistency model is a framework to test if a particular system is comply with a specific consistency model. 

## Strong consistency model
- Linearizability: the strongest consistency model for single operation. It is atomic
  and guarantee it is consistent with the real time ordering of that event.
- Serializabitily: consistency model for transaction operation. It doesn't guarantee that the real time order is the same as the execution order. Serializability guarantees that operations take place atomically: a transaction’s sub-operations do not appear to interleave with sub-operations from other transactions.

    **Atomicity** : Each transaction is all-or-nothing; sub-operations do not interleave in a way that violates serializability.

    **No real-time constraints**: Even though T1 started first, T2 could be “seen” as executed first.

    **Multi-object** : The transaction T1 updates A and B, yet the system maintains a consistent view.

- Strict Serializability: combines the linearizability and the serializability where 
  the concurrent process will be done in some order and the real time order is the same 
  with the system's execution order. Strict serializability ensures log order == real-time order.

## Weaker consistency model
### 1. **Sequential Consistency**

**Definition:** Operations appear to take place in some total order, and that order is consistent with the order of operations on each individual process.

**Analogy:** Imagine a line of people passing a ball. Each person passes the ball in the order they receive it, and everyone observes the same sequence of passes.

**Real-World Example:** In a collaborative text editor, if one user types a character, all other users will see that character appear in the same position, in the same order.

---

### 2. **Causal Consistency (Monotonic Reads + Monotonic Writes + Read Your Writes + Writes Follows Read)**

**Definition:** Operations that are causally related must be seen by all processes in the same order — though processes may disagree about the order of causally independent operations

**Analogy:** If Person A tells Person B a secret, and Person B tells Person C, then Person C knows the secret because it's causally related.

**Real-World Example:** In a social media application, if User A posts a comment and User B replies to it, all users will see the reply after the comment, preserving the causal relationship.

---

### 3. **PRAM (Monotonic Reads + Monotonic Writes + Read Your Writes)**

**Definition:** Combines three properties:

* **Monotonic Reads:** If a process reads a value, any subsequent reads will return the same or a more recent value.
* **Monotonic Writes:** Writes by a process are serialized in the order they were issued.
* **Read Your Writes:** A process will always see its own writes.

**Analogy:** If you place a book on a shelf and then move it, you'll always see the book in its new position.

**Real-World Example:** In an online shopping cart, if you add an item, remove it, and then add it again, the cart will consistently reflect the most recent actions.

---

### 4. **Read Follows Write**

**Definition:** If a process writes a value and then reads the same value, it must see the value it just wrote.

**Analogy:** If you write a note and then read it, you'll see the exact words you wrote.

**Real-World Example:** In a banking application, after transferring money, if you check your balance, it should reflect the new amount immediately.

---

### 5. **Monotonic Read**

**Definition:** If a process reads a value, any subsequent reads will return the same or a more recent value.

**Analogy:** If you check the weather forecast and then check it again later, the forecast will be the same or more up-to-date.

**Real-World Example:** In a news application, once you read an article, any future reads will show the same or newer content.

---

### 6. **Monotonic Write**

**Definition:** Writes by a process are serialized in the order they were issued.

**Analogy:** If you send two emails, the second email will be sent after the first, and recipients will see them in that order.

**Real-World Example:** In a document editing application, if you save changes and then make more changes, the second save will reflect the changes made after the first.

---

### 7. **Read Your Own Write**

**Definition:** A process will always see its own writes.

**Analogy:** If you write a message on a whiteboard and then read it, you'll see the exact message you wrote.

**Real-World Example:** In a messaging application, after sending a message, you can immediately see it in the chat history.

---

## Notes
- Strong consistency => a process will be consistent accross all nodes. Not available during fault because of it have to consistent. Also complex to implement. 
- Eventual consistency => process that will eventually be consistent (convergence) given a specific amount of time, which can't be known exactly. This is a weaker consistency, but available during fault. 
- Causal consistency => the order of any causal processes will be consistent across nodes even though the independent process may not have the same order

** Example of system behaviour under strong consistency **
- Insert data to db, all db replica has the most updated value immediately
- Do some bank transfer, A:x=>=>y=>z and B=>j=>k=>l, both transaction is guaranteed to be atomic, so there will be no interleaving for y or k, or any suboperations. But the order of A after B or B after
A is not guaranteed.

** Example of system behaviour under weak consistency **
- I post a comment, I can read the most recent comment I've made, but another user won't see it immediately. (read my write)
- I reply a invitation chat in a group, but the order of another members' chat after that invitation might not be the same as the real time order. (causal consistency)
- I read by bank account ow, if I read the account later on and the value is updated, it has been guaranteed that the latter is the most recent data (monotonic read) 