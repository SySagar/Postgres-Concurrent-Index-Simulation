#include <iostream>
#include <vector>
#include <string>
#include <map>

#include <thread>
#include <chrono>
#include <mutex>

using namespace std;

struct User{
    int  id;
    string name;
};

mutex mtx;

vector<User> table;
map<string,int> idx;   // fake index

vector<User> pending;
bool building = false;

void synchronize() {

    lock_guard<mutex> lock(mtx);

    cout << "Synchronizing..." << endl;

    for (auto user : pending) {
        idx[user.name] = user.id;
    }

    pending.clear();
}


void buildIndex(){

    lock_guard<mutex> lock(mtx); //builder holds the mutex for the entire index build
     building = true;

    for(auto user: table){
        idx[user.name] = user.id;

          this_thread::sleep_for(chrono::seconds(1)); // simulates millions row transaction scan and indexing - 1 s gap each row
    }

    synchronize();

    building = false;
}

void insertUser(int id, string name) {

     lock_guard<mutex> lock(mtx);  // lock mutex -> safely execture -> unlocks

     User user{id, name};

    table.push_back(user);

    cout << "Inserted: " << name << endl;

    if (building) {
        pending.push_back(user); // remember it as it wasn't caught during scanning
    }
}

int main() {
    
    table.push_back({1, "Alice"});
    table.push_back({2, "Bob"});
    table.push_back({3, "Charlie"});
    table.push_back({3, "Charlie"});
    table.push_back({4, "David"});
    table.push_back({5, "Eve"});

    thread builder(buildIndex);

    this_thread::sleep_for(chrono::seconds(2));

      insertUser(6, "Frank");

      builder.join();


   for (auto entry : idx) {
        cout << entry.first
             << " -> "
             << entry.second
             << endl;
    }

    return 0;
}