# Programming Paradigm - Shipping Yard #

I solved this question using Java programming language. It's because in my opinion, Java is the best when it comes to Object-Oriented Programming Paradigm.

I tested using this code:
```
public static void main(String args[]) {
    SeaVehicle yacht = new Sailboat("My Yacht", "BLUEPRINT OF YACHT");

    System.out.println("This is a " + yacht.getType());
    System.out.println("Its name is " + yacht.getName());
    System.out.println("This is how it looks: " + yacht.getBlueprint());

    yacht.setName("Not My Yacht anymore");
    System.out.println("I have changed its name to " + yacht.getName());
}
```
and it gives output:
```
This is a sailboat
Its name is My Yacht
This is how it looks: BLUEPRINT OF YACHT
I have changed its name to Not My Yacht anymore
```

See original question for this problem below.

## [ORIGINAL QUESTION] ##

Sebuah galangan kapal ingin membuat aplikasi untuk menyimpan data/denah kapal yang telah dibuat.

Dengan konsep pemrograman berorientasi objek, buatlah class beserta abstraksi, function dan properties untuk beberapa tipe kapal(perahu motor, perahu layar & kapal pesiar) dan mengimplementasikan beberapa konsep dasar seperti enkapsulasi & polimorfisme.

Class harap di upload ke PRIVATE GitHub repository dan tambahkan/invite fandywie sebagai Collaborator.
Apabila ada pertanyaan dapat menghubungi tech.career@sirclo.co.id 