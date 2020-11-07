abstract class SeaVehicle {
    private String name;
    private String blueprint;

    public SeaVehicle (String name, String blueprint) {
        this.name = name;
        this.blueprint = blueprint;
    }
    
    public abstract String getType();

    public String getName() {
        return this.name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getBlueprint() {
        return this.blueprint;
    }

    public void setBlueprint(String blueprint) {
        this.blueprint = blueprint;
    }
}

class Motorboat extends SeaVehicle {
    public Motorboat(String name, String blueprint) {
        super(name, blueprint);
    }

    public String getType() {
        return "motorboat";
    }
}

class Sailboat extends SeaVehicle {
    public Sailboat(String name, String blueprint) {
        super(name, blueprint);
    }

    public String getType() {
        return "sailboat";
    }
}

class Cruise extends SeaVehicle {
    public Cruise(String name, String blueprint) {
        super(name, blueprint);
    }

    public String getType() {
        return "cruise";
    }

}