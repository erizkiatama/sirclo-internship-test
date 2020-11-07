import java.util.Map;
import java.util.TreeMap;

public class Cart {
    private final Map<String, Integer> keranjang = new TreeMap<>();

    void tambahProduk(String kodeProduk, int kuantitas) {
        if (!keranjang.containsKey(kodeProduk)) {
            keranjang.put(kodeProduk, kuantitas);
        } else {
            int jumlahSekarang = keranjang.get(kodeProduk);
            keranjang.replace(kodeProduk, jumlahSekarang + kuantitas);
        }
    }

    void hapusProduk(String kodeProduk) {
        keranjang.remove(kodeProduk);
    } 

    void tampilkanCart() {
        if (keranjang.isEmpty()) {
            System.out.println("Keranjang belanja anda kosong");
        } else {
            for (String kodeProduk : keranjang.keySet()) {
                System.out.printf("%s (%d)\n", kodeProduk, keranjang.get(kodeProduk));
            }
        }
    }

    public Map<String, Integer> getKeranjang() {
        return keranjang;
    }
}
