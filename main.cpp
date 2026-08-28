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

void buildIndex(){
    for(auto user: table){
        idx[user.name] = user.id;

          this_thread::sleep_for(chrono::seconds(1)); // simulates millions row transaction scan and indexing - 1 s gap each row
    }
}

void insertUser(int id, string name) {

    table.push_back({id, name});

    cout << "Inserted: " << name << endl;
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