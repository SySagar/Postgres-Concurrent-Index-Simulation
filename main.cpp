#include <iostream>
#include <vector>
#include <string>
#include <map>

#include <thread>
#include <chrono>

using namespace std;

struct User{
    int  id;
    string name;
};

vector<User> table;
map<string,int> idx;   // fake index

vector<User> pending;
bool building = false;

void synchronize() {

    cout << "Synchronizing..." << endl;

    for (auto user : pending) {
        idx[user.name] = user.id;
    }

    pending.clear();
}


void buildIndex(){
     building = true;

    for(auto user: table){
        idx[user.name] = user.id;

          this_thread::sleep_for(chrono::seconds(1)); // simulates millions row transaction scan and indexing - 1 s gap each row
    }

    synchronize();

    building = false;
}

void insertUser(int id, string name) {

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